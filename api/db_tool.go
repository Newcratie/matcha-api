package api

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

func customQuery(Id int, Filter *Filters) (superQuery string) {
	var cQuery string

	minAge := ageConvert(Filter.Age[0])
	maxAge := ageConvert(Filter.Age[1])

	if len(Filter.Tags) > 0 {
		cQuery = setTagQuery(Filter)
	}

	superQuery += `MATCH (u:User) WHERE NOT Id(u)= ` + strconv.Itoa(Id) + ` AND (u.rating >= ` + strconv.Itoa(Filter.Score[0]) + ` AND u.rating <= ` + strconv.Itoa(Filter.Score[1]) + `)
	AND (u.birthday >= "` + maxAge + `" AND u.birthday <= "` + minAge + `") ` + cQuery + ` RETURN DISTINCT u`

	return
}

func setTagQuery(Filter *Filters) (customQuery string) {

	customQuery = `MATCH (t:TAG)-[]-(u) WHERE `
	for i, tag := range Filter.Tags {
		if i == 0 {
			customQuery += `t.value='` + tag + `' `
		} else {
			customQuery += `OR t.value='` + tag + `' `
		}
	}
	return
}

//func setContent()

func ageConvert(Age int) (birthYear string) {

	now := time.Now()
	now = now.AddDate(-(Age), 0, 0)
	birthYear = now.Format(time.RFC3339Nano)
	return birthYear
}

func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance int) {

	const earthRadius = float64(6371)

	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = int(earthRadius * c)

	return
}

var floatType = reflect.TypeOf(float64(0))

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

// Create link between User and TAG
//MATCH (n:User {name: 'Yundt7209'}) MATCH (t:TAG {name: 'css'}) CREATE (n)-[:TAGGED]->(t)

// Delete link between User and TAG
//MATCH (n:User {name: 'Yundt7209'}) MATCH (t:TAG {name: 'css'}) MATCH (n)-[l:TAGGED]-(t) DELETE l

// Sample of a query with tag
//MATCH (u:User) WHERE (u.rating > 0 AND u.rating < 51) MATCH (u) WHERE (u.birthday > "1899-04-18T17:18:43.342718527Z" AND u.birthday < "2003-04-18T17:18:43.342712091Z") MATCH (t:TAG {name: 'css'})-[]-(u) RETURN u

//MATCH (t:TAG)-[]-(u) WHERE t.name='design' OR t.name='HTML' RETURN u
