package database

import "todo-api-go/prisma/db"

func CreateConnection() *db.PrismaClient {
	prismaClient := db.NewClient()
	prismaClient.Connect()

	return prismaClient
}

var PrismaClient = CreateConnection()
