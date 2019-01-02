// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package man provides a concurrency-safe bank with one account.
package main

import (
	"fmt"
	"time"
)

type withdrawMsg struct {
	amount int
	status chan bool
}

var deposits = make(chan int)          // send amount to deposit
var balances = make(chan int)          // receive balance
var withdraws = make(chan withdrawMsg) // send amount to deposit, receive withdraw status

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	msg := withdrawMsg{amount, make(chan bool)}
	withdraws <- msg
	return <-msg.status
}

func monitor() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case msg := <-withdraws:
			if msg.amount > balance {
				msg.status <- false
			} else {
				balance -= msg.amount
				msg.status <- true
			}
		}
	}
}

func main() {
	done := make(chan struct{})

	// Alice
	go func() {
		Deposit(200)
		fmt.Println("deposit =", Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		Deposit(100)
		fmt.Println("deposit =", Balance())
		done <- struct{}{}
	}()

	go func() {
		time.Sleep(300 * time.Millisecond)
		status := Withdraw(500)
		fmt.Printf("withdraw ")
		if status {
			fmt.Println("succeed")
		} else {
			fmt.Println("failed")
		}
		done <- struct{}{}
	}()

	// Wait for all transactions.
	<-done
	<-done
	<-done
}

func init() {
	go monitor() // start the monitor goroutine
}

//!-
