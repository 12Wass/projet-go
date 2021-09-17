package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Player struct {
	id    int
	name  string
	boats [10][10]int
}

var playerOne Player // Joueur 1
var playerTwo Player // Joueur 2

func helloHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Hello world")
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
	}
}

func dateHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		currentTime := time.Now()
		fmt.Fprintf(w, "%s", currentTime.Format("03h04"))
	}
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if err := req.ParseForm(); err != nil { // Parsing des paramètres envoyés
			fmt.Println("Something went bad") // par le client et gestion d’erreurs
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm { // On print les clés et valeurs des
			fmt.Println(key, "=>", value) // données envoyés par le clients
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
		// Sauvegarde des données reçues
		saveFile, err := os.OpenFile("./save.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		defer saveFile.Close()

		w := bufio.NewWriter(saveFile)
		if err == nil {
			fmt.Fprintf(w, "%v:%v:\n", req.PostForm["entries"][0], req.PostForm["author"][0])
		}
		w.Flush()

	}
}

func entriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		filerc, err := os.Open("./save.data")
		if err != nil {
			log.Fatal(err)
		}
		defer filerc.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(filerc)
		contents := buf.String()

		split := strings.Split(contents, ":")

		for k := range split {
			if k%2 == 0 {
				fmt.Fprintf(w, split[k])
			}
		}
	}
}

func initGame() {

	playerOne.id = 1
	playerTwo.id = 2

	fmt.Println("Bienvenue dans ce jeu, veuillez saisir le nom du premier joueur : ")
	fmt.Scan(&playerOne.name)
	initBoat(&playerOne)

	fmt.Println("Second joueur : ")
	fmt.Scan(&playerTwo.name)
	initBoat(&playerTwo)

}
func initBoat(player *Player) {
	boatCount := 3
	fmt.Println("Pour placer vos bateaux, écrivez à la suite les coordonnées x et y ainsi que la direction souhaitée (up, down, left, right)")
	var xPosition int
	var yPosition int

	for boatCount > 0 {
		fmt.Println("Vous pouvez placer", boatCount, "bateaux")

		fmt.Print("X : ")
		fmt.Scan(&xPosition)
		fmt.Println(xPosition)

		fmt.Print("Y : ")
		fmt.Scan(&yPosition)
		fmt.Println(yPosition)
		if !checkPosition(xPosition, yPosition, player) || xPosition > 10 || yPosition > 10 {
			player.boats[xPosition][yPosition] = 1
			boatCount--
		} else if xPosition > 10 || yPosition > 10 {
			fmt.Println("Veuillez réessayer, les positions dépassent du tableau")
		} else {
			fmt.Println("Veuillez réessayer, les positions sont déjà prises")
		}
	}
}

func checkPosition(x int, y int, player *Player) bool {
	if player.boats[x][y] > 0 {
		return true
	} else {
		return false
	}
}

func boardHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Joueur 1 : %s \n", playerOne.name)
		for _, j := range playerOne.boats {
			fmt.Fprintf(w, "%v\n", j)
		}
		fmt.Fprintf(w, "Joueur 2 : %s \n", playerTwo.name)
		for _, j := range playerTwo.boats {
			fmt.Fprintf(w, "%v\n", j)
		}
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", dateHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/entries", entriesHandler)
	http.HandleFunc("/board", boardHandler)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":4567", nil))
		wg.Done()
	}()
	initGame()
	wg.Wait()
}
