package handler_auth

import (
	"context"
	"log"
	pbTimekeeping "skko-timekeeping-ms/pb"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func Timekeeping(c *fiber.Ctx) error {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Println("❌ failed to connect gRPC:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect gRPC"})
	}
	defer conn.Close()

	client := pbTimekeeping.NewTimekeepingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetInitialData(ctx, &pbTimekeeping.InitialDataRequest{
		UserId: 123,
	})
	if err != nil {
		log.Println("❌ gRPC error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gRPC request failed"})
	}

	return c.JSON(resp)
}
