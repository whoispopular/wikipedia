WIKIPEDIA
=========

Set of tools to deal with wikipedia dumps, filter, transform, insert to database
and serve wikipedia data.

# Snippets

To generate people XML database from online wikipedia dump

```
export WIKILANG=ar
curl https://dumps.wikimedia.org/${WIKILANG}wiki/latest/${WIKILANG}wiki-latest-pages-articles-multistream.xml.bz2 | go run cmd/wikipedia-extract/wikipedia-extract.go -config configs/people/${WIKILANG}.json
```

To insert pages to people table

```
cat ~/Downloads/wikipedia/people.ar.xml | go run cmd/wikipedia-insert/wikipedia-insert.go --language=ar --entity=Person --table=people
```
