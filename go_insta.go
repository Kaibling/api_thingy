package main

import (
	"github.com/ahmdrz/goinsta"
	"fmt"
	"log"
)


func fetchTag(insta *goinsta.Instagram, tag string) error {

	feedTag, err := insta.Feed.Tags(tag)
	if err != nil {
		return err
	}
	for _, item := range feedTag.RankedItems {
		log.Printf("error on liking item %s, %v", item.ID, err)

	}
	return nil
}

func main() {
	insta := goinsta.New(
		username,
		password,
	)
	if err := insta.Login(); err != nil {
		log.Println(err)
		return
	}
	defer insta.Logout()

	a,err := insta.Profiles.ByName("asads")



	if err != nil {
		log.Println(err)
	}
    log.Printf("%s %d", a.FullName, a.ID)
    amount := a.Followers().Users
    log.Println(amount)
	feedi := a.Feed()
    log.Println(feedi.Items)
	log.Println(a.Stories().Items)

	for feedi.Next(false) {
		for _, item := range feedi.Items {
log.Printf("id: %s\n",item.ID)
		}
	}

}