package main

import "db_labs/ioc"

func main() {
	print(ioc.UsePgConnection())
}