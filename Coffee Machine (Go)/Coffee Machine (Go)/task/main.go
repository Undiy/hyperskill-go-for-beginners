package main

import "fmt"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printStatus(water, milk, beans, cups, money, waste int) {
	fmt.Printf(
		`The coffee machine has:
%d ml of water
%d ml of milk
%d g of coffee beans
%d disposable cups
%d of money
%d g of waste
`, water, milk, beans, cups, money, waste)
}

func buyCoffee(water, milk, beans, cups, money, waste *int, w, m, b, mm int) {
	if *water < w {
		fmt.Println("Sorry, not enough water!")
	} else if *milk < m {
		fmt.Println("Sorry, not enough milk!")
	} else if *beans < b {
		fmt.Println("Sorry, not enough coffee beans!")
	} else if *cups < 1 {
		fmt.Println("Sorry, not enough disposable cups!")
	} else if *waste > 100 {
		fmt.Println("Sorry, the coffee machine should be cleaned first!")
	} else {
		fmt.Println("I have enough resources, making you a coffee!")
		*water -= w
		*milk -= m
		*beans -= b
		*cups -= 1
		*money += mm
		*waste += b
	}
}

func buyAction(water, milk, beans, cups, money, waste *int) {
	fmt.Println("What do you want to buy? 1 - espresso, 2 - latte, 3 - cappuccino, back - to main menu: ")
	var choice string
	fmt.Scan(&choice)
	switch choice {
	case "1":
		buyCoffee(
			water,
			milk,
			beans,
			cups,
			money,
			waste,
			250,
			0,
			16,
			4,
		)
	case "2":
		buyCoffee(
			water,
			milk,
			beans,
			cups,
			money,
			waste,
			350,
			75,
			20,
			7,
		)
	case "3":
		buyCoffee(
			water,
			milk,
			beans,
			cups,
			money,
			waste,
			200,
			100,
			12,
			6,
		)
	}
}

func fillAction(water, milk, beans, cups *int) {
	var w, m, b, c int

	fmt.Println("Write how many ml of water you want to add:")
	fmt.Scan(&w)

	fmt.Println("Write how many ml of milk you want to add:")
	fmt.Scan(&m)

	fmt.Println("Write how many grams of coffee beans you want to add:")
	fmt.Scan(&b)

	fmt.Println("Write how many disposable cups you want to add:")
	fmt.Scan(&c)

	*water += w
	*milk += m
	*beans += b
	*cups += c
}

func takeAction(money *int) {
	fmt.Printf("I gave you $%d\n", *money)
	*money = 0
}

func cleanAction(waste *int) {
	fmt.Println("I'm clean now!")
	*waste = 0
}

func main() {
	water, milk, beans, cups, money := 400, 540, 120, 9, 550
	waste := 0

	var action string
	for action != "exit" {
		fmt.Println("Write action (buy, fill, take, clean, remaining, exit):")
		fmt.Scan(&action)
		switch action {
		case "buy":
			buyAction(&water, &milk, &beans, &cups, &money, &waste)
		case "fill":
			fillAction(&water, &milk, &beans, &cups)
		case "take":
			takeAction(&money)
		case "clean":
			cleanAction(&waste)
		case "remaining":
			printStatus(water, milk, beans, cups, money, waste)
		}
		fmt.Println("")
	}
}
