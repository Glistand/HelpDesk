package main

import (
	"fmt"
	"ticket-service/internal/store"
)

func main() {
	s := store.NewStore()

	t1, err := s.Create("VPN не подключается")
	if err != nil {
		fmt.Println("create t1:", err)
		return
	}
	fmt.Println("created:", t1.ID, t1.Title)

	t2, err := s.Create("Не работает почта")
	if err != nil {
		fmt.Println("create t2:", err)
		return
	}
	fmt.Println("created:", t2.ID, t2.Title)

	fmt.Println("--- list ---")
	for _, t := range s.List() {
		fmt.Println(t.ID, t.Status, t.Title)
	}

	fmt.Println("--- get ---")
	got, err := s.Get(t1.ID)
	if err != nil {
		fmt.Println("get:", err)
		return
	}
	fmt.Println("got:", got.Title)

	_, err = s.Create("")
	fmt.Println("empty title err:", err)

	_, err = s.Get("t-999")
	fmt.Println("not found err:", err)
}
