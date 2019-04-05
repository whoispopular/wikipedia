package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"wikipedia"
)

func main() {
	lang := flag.String("language", "en", "language code, will be used for the name column eg. ar will use name_ar")
	table := flag.String("table", "people", "the entities table to insert entities to")
	entity := flag.String("entity", "Person", "entity type")
	flag.Parse()
	db_url := os.Getenv("DATABASE_URL")

	parser, err := wikipedia.NewParser(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	pages := make(chan wikipedia.Page, 1000)
	go wikipedia.ReadAllPages(*parser, pages)
	insertPages(pages, db_url, *table, *lang, *entity)
}

func insertPages(pages <-chan wikipedia.Page, db_url string, table string, lang string, entity string) {
	db := wikipedia.OpenDB(db_url)
	defer wikipedia.CloseDB(db)

	for page := range pages {
		if strings.Index(page.Title, " ") == -1 {
			continue
		}

		entityId, err := getOrInsertEntity(page, db, table, lang, entity)
		if err != nil {
			log.Printf("%s, %s", page.Title, err)
			continue
		}

		_, err = insertOrUpdatePage(page, db, lang, entity, entityId)
		if err != nil {
			log.Printf("%s, %s", page.Title, err)
		}
	}
}

func getOrInsertEntity(page wikipedia.Page, db wikipedia.DB, table string, lang string, entity string) (id int, err error) {
	if isEntityExists(page.Title, db, table, lang) {
		err = db.QueryRow("SELECT id FROM "+table+" WHERE name_"+lang+" = $1", page.Title).Scan(&id)
		return
	}

	return insertEntity(page, db, table, lang, entity)
}

func isEntityExists(title string, db wikipedia.DB, table string, lang string) bool {
	var count int
	db.QueryRow("SELECT COUNT(id) as c FROM "+table+" WHERE name_"+lang+" = $1", title).Scan(&count)
	return count > 0
}

func insertEntity(page wikipedia.Page, db wikipedia.DB, table string, lang string, entity string) (id int, err error) {
	err = db.QueryRow(`INSERT INTO `+table+`(name_`+lang+`, created_at, updated_at)
					VALUES($1, NOW(), NOW())
					RETURNING id
          `, page.Title).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("Error inserting person: %s", err)
	}

	var synonymid int
	err = db.QueryRow(`INSERT INTO synonyms(entity_id, entity_type, locale, synonym, created_at, updated_at)
					VALUES($1, $2, $3, $4, NOW(), NOW())
					RETURNING id
          `, id, entity, lang, page.Title).Scan(&synonymid)

	if err != nil {
		return id, fmt.Errorf("Error inserting synonym: %s", err)
	}

	return
}

func insertOrUpdatePage(page wikipedia.Page, db wikipedia.DB, lang string, entity string, entityId int) (id int, err error) {
	if isPageExists(int(page.ID), db, lang) {
		return 0, updatePage(page, db, entityId, entity, lang)
	}

	return insertPage(page, db, entityId, entity, lang)
}

func isPageExists(pageCurId int, db wikipedia.DB, lang string) bool {
	var count int
	db.QueryRow("SELECT COUNT(id) as c FROM wikipedia_pages WHERE curid = $1 AND locale = $2", pageCurId, lang).Scan(&count)
	return count > 0
}

func insertPage(page wikipedia.Page, db wikipedia.DB, entityId int, entity string, lang string) (id int, err error) {
	err = db.QueryRow(`INSERT INTO wikipedia_pages(curid, excerpt, entity_id, entity_type, locale, created_at, updated_at)
					VALUES($1, $2, $3, $4, $5, NOW(), NOW())
					RETURNING id
          `, page.ID, wikipedia.FirstParagraph(page.Text), entityId, entity, lang).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("Error inserting wikipedia page: %s", err)
	}

	return
}

func updatePage(page wikipedia.Page, db wikipedia.DB, entityId int, entity string, lang string) (err error) {
	_, err = db.Exec(`UPDATE wikipedia_pages
          SET excerpt = $1, entity_id = $2, entity_type = $3, updated_at = NOW()
					WHERE locale = $4 AND curid = $5
          `, wikipedia.FirstParagraph(page.Text), entityId, entity, lang, page.ID)

	if err != nil {
		return fmt.Errorf("Error updating wikipedia page: %s", err)
	}

	return
}
