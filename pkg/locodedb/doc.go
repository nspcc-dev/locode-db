/*
Package locodedb implements a UN LOCODE database.

It contains all the data internally and provides simple [Get] API to retrieve
records based on short LOCODE strings. The DB is stored compressed before the
first use (~1MB) and is unpacked automatically on the first access (which takes
~100-200ms). Unpacked it needs ~4MB of RAM.
*/
package locodedb
