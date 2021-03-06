package main

import (
	"fmt"
	"strings"
	"unicode"
	"strconv"
	"github.com/pauloaguiar/ces27-lab1/mapreduce"
	"hash/fnv"
	//"reflect"
)

// mapFunc is called for each array of bytes read from the splitted files. For wordcount
// it should convert it into an array and parses it into an array of KeyValue that have
// all the words in the input.
func mapFunc(input []byte) (result []mapreduce.KeyValue) {
	// 	Pay attention! We are getting an array of bytes. Cast it to string.
	//
	// 	To decide if a character is a delimiter of a word, use the following check:
	//		!unicode.IsLetter(c) && !unicode.IsNumber(c)
	//
	//	Map should also make words lower cased:
	//		strings.ToLower(string)
	//
	// IMPORTANT! The cast 'string(5)' won't return the character '5'.
	// 		If you want to convert to and from string types, use the package 'strconv':
	// 			strconv.Itoa(5) // = "5"
	//			strconv.Atoi("5") // = 5

	/////////////////////////
	// YOUR CODE GOES HERE //
	/////////////////////////

	str_in := string(input[:])
	str_in = str_in + " "
	str_in = strings.ToLower(str_in)

	myMap := make(map[string]int)
	word := ""

	for _, c := range str_in {
		if  unicode.IsLetter(c) || unicode.IsNumber(c) {
			word += string(c)
		} else {
			if word != "" {
				if _, ok := myMap[word]; !ok{
					myMap[word] = 1
				} else {
					myMap[word]++
				}
			}
			word = ""
		}
	}
	
	for key, value := range myMap {
		//fmt.Printf("key: %v\n", key)
		result = append(result, mapreduce.KeyValue{key, strconv.Itoa(value)})
	}

	fmt.Printf("result length: %v\n", len(result))

	// for i:=0; i<len(result); i++{
	// 	fmt.Printf("key: %v, value: %v\n", result[i].Key, result[i].Value)
	// }

	//result = make([]mapreduce.KeyValue, 0)
	return result
}

// reduceFunc is called for each merged array of KeyValue resulted from all map jobs.
// It should return a similar array that summarizes all similar keys in the input.
func reduceFunc(input []mapreduce.KeyValue) (result []mapreduce.KeyValue) {
	// 	Maybe it's easier if we have an auxiliary structure? Which one?
	//
	// 	You can check if a map have a key as following:
	// 		if _, ok := myMap[myKey]; !ok {
	//			// Don't have the key
	//		}
	//
	// 	Reduce will receive KeyValue pairs that have string values, you may need
	// 	convert those values to int before being able to use it in operations.
	//  	package strconv: func Atoi(s string) (int, error)
	//
	// 	It's also possible to receive a non-numeric value (i.e. "+"). You can check the
	// 	error returned by Atoi and if it's not 'nil', use 1 as the value.

	/////////////////////////
	// YOUR CODE GOES HERE //
	/////////////////////////

	myMap := make(map[string]int)

	for i := 0; i < len(input); i++ {
		value, err := strconv.Atoi(input[i].Value)
		if err != nil {
			value = 1
		}
		if _, ok := myMap[input[i].Key]; !ok {
			myMap[input[i].Key] = value
		} else {
			myMap[input[i].Key] += value
		}
	}

	for key, value := range myMap {
		result = append(result, mapreduce.KeyValue{key, strconv.Itoa(value)})
	}

	// for i:=0; i<len(result); i++{
	// 	fmt.Printf("key: %v, value: %v\n", result[i].Key, result[i].Value)
	// }

	//result = make([]mapreduce.KeyValue, 0)
	return result
}

// shuffleFunc will shuffle map job results into different job tasks. It should assert that
// the related keys will be sent to the same job, thus it will hash the key (a word) and assert
// that the same hash always goes to the same reduce job.
// http://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
func shuffleFunc(task *mapreduce.Task, key string) (reduceJob int) {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() % uint32(task.NumReduceJobs))
}
