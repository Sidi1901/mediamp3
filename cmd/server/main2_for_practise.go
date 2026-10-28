package main

import(
	"fmt"
)

func printTags(args ...string){
	for i, tag := range args{
		fmt.Printf("This is the index %d and this is word %s \n", i, tag)
	}
}

func main(){
	myslice := []string{"hello","sidi", "Bhardwaj", "."}
	printTags(myslice...)
}