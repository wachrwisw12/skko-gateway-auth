package handlerleave

import (
	"context"
	"log"
	pbLeave "skko-leave-ms/pb"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type RequestBody struct {
	UserId int `json:"user_id"`
	SexId  int `json:"sex_id"`
}

func GetleaveHandler(c *fiber.Ctx) error {
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}
	log.Print(body.SexId)
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Println("❌ failed to connect gRPC:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect gRPC"})
	}

	defer conn.Close()

	client := pbLeave.NewLeaveServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.GetLeave(ctx, &pbLeave.GetLeaveRequest{
		UserId: int32(body.UserId),
		SexId:  int32(body.SexId),
	})
	if err != nil {
		log.Println("❌ gRPC error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gRPC request failed"})
	}

	return c.JSON(resp)
}
