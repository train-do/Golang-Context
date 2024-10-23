package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"
)
type User struct{
	username string
	keranjang [][]string
}
func main() {
	var users []User
	users = append(users, User{"user1", [][]string{{"Handphone", "1", "$20"}}})
	users = append(users, User{"user2", [][]string{{"Laptop", "1", "$20"}, {"Handphone", "1", "$20"}}})
	for{
		var input string
		fmt.Println("--- Halaman Login ---")
		fmt.Print("Masukkan Username : ")
		fmt.Scanln(&input)
		user, isValid, close := validasi(input, users)
		if isValid != nil {
			ClearScreen()
			fmt.Println(isValid)
			continue
		}
		if close {
			os.Exit(0)
		}
		ctx, cancel :=  context.WithTimeout(context.Background(), 20 * time.Second)
		runningApp(ctx, user, cancel)
		select {
		case <- ctx.Done():
			if reflect.TypeOf(ctx.Err()).Name() == "deadlineExceededError"{
				fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m\n", ctx.Err())
			}else{
				fmt.Printf("Anda Berhasil Logout \033[31m%v\033[0m\n", ctx.Err())
			}
			cancel()
		}
	}
}
func runningApp(ctx context.Context, user User, cancel context.CancelFunc)  {
	for{
		select {
		case <- ctx.Done():
			return
		default:
			ClearScreen()
			printMenu()
			var input string
			fmt.Scanln(&input)
			switch input {
			case "1":
				printProduk()
			case "2":
				printKeranjang(user.keranjang)
			case "3":
				printCheckout(user.keranjang)
			case "0":
				ClearScreen()
				cancel()
			}
		}
	}	
}
func validasi(input string, users []User) (User, error, bool) {
	if input == "exit" {
		return User{}, nil, true
	}
	for _, v := range users {
		if input == v.username {
			return v, nil, false
		}
	}
	return User{}, errors.New("Invalid Username"), false
}
func printMenu()  {
	fmt.Println("--- MAIN MENU ---")
	fmt.Println("1. Show Produk")
	fmt.Println("2. Show Keranjang")
	fmt.Println("3. Checkout")
	fmt.Println("0. Logout")
}
func printProduk()  {
	ClearScreen()
	fmt.Println("--- Halaman Produk ---")
	fmt.Println("1. Laptop")
	fmt.Println("2. Handphone")
	fmt.Println("3. PC")
	fmt.Scanln()
}
func printKeranjang(keranjang [][]string)  {
	ClearScreen()
	fmt.Println("--- Halaman Keranjang ---")
	for i, v := range keranjang {
		fmt.Printf("%d. %s\t\t%spcs\n", (i+1), v[0], v[1])
	}
	fmt.Scanln()
}
func printCheckout(keranjang [][]string)  {
	ClearScreen()
	fmt.Println("--- Halaman Checkout ---")
	for i, v := range keranjang {
		fmt.Printf("%d. %s\t\t%spcs\t\t%s\n", (i+1), v[0], v[1], v[2])
	}
	fmt.Scanln()
}
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}