package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowWidth  = 512
	windowHeight = 380

	numberOfEnemies = 4

	gravity      = 1000
	jumpVelocity = -600

	safetyPad = 40
)

type animData struct {
	rec         rl.Rectangle
	pos         rl.Vector2
	frame       int
	maxFrame    int
	updateTime  float32
	runningTime float32
	velocity    float32
}

func (a *animData) updateAnimData(deltaTime float32) {
	a.runningTime += deltaTime
	if a.runningTime >= a.updateTime {
		a.runningTime = 0.0
		a.rec.X = float32(a.frame) * a.rec.Width

		a.frame++
		if a.frame > a.maxFrame {
			a.frame = 0
		}
	}
}

func isOnGround(player animData, windowHeight int) bool {
	return player.pos.Y >= float32(windowHeight)-player.rec.Height
}

func main() {
	rl.InitWindow(windowWidth, windowHeight, "Dapper Dasher")
	rl.SetTargetFPS(60)

	// start - nebula
	nebulaTexture := rl.LoadTexture("./textures/12_nebula_spritesheet.png")
	enemiesNebula := make([]animData, 0, numberOfEnemies)

	for i := 0; i < numberOfEnemies; i++ {
		enemiesNebula = append(enemiesNebula, animData{
			rec:         rl.NewRectangle(0.0, 0.0, float32(nebulaTexture.Height)/8, float32(nebulaTexture.Width)/8),
			pos:         rl.NewVector2(float32(windowWidth)+300.0*float32(i), windowHeight-float32(nebulaTexture.Height)/8),
			frame:       0,
			maxFrame:    7,
			updateTime:  1.0 / 12.0,
			runningTime: 0.0,
			velocity:    -200,
		})
	}
	// end - nebula

	// start - scarfy
	scarfyTexture := rl.LoadTexture("./textures/scarfy.png")

	scarfy := animData{
		rec:         rl.NewRectangle(0.0, 0.0, float32(scarfyTexture.Width)/6, float32(scarfyTexture.Height)),
		pos:         rl.NewVector2(windowWidth/2-float32(scarfyTexture.Width)/12, windowHeight-float32(scarfyTexture.Height)),
		frame:       0,
		maxFrame:    5,
		updateTime:  1.0 / 12.0,
		runningTime: 0.0,
		velocity:    0,
	}

	scarfyIsInAir := false
	// end - scarfy

	// start - background
	backgroundTexture := rl.LoadTexture("./textures/far-buildings.png")
	midgroundTexture := rl.LoadTexture("./textures/back-buildings.png")
	foregroundTexture := rl.LoadTexture("./textures/foreground.png")

	backgroundPosition := rl.NewVector2(0, 0)
	background2Position := rl.NewVector2(float32(backgroundTexture.Width)*2, 0)

	midgroundPosition := rl.NewVector2(0, 0)
	midground2Position := rl.NewVector2(float32(midgroundTexture.Width)*2, 0)

	foregroundPosition := rl.NewVector2(0, 0)
	foreground2Position := rl.NewVector2(float32(foregroundTexture.Width)*2, 0)
	// end - background

	// start - conditions
	collision := false
	finishLine := enemiesNebula[len(enemiesNebula)-1].pos.X
	// end - conditions

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		backgroundPosition.X -= 20 * deltaTime
		background2Position.X -= 20 * deltaTime

		midgroundPosition.X -= 40 * deltaTime
		midground2Position.X -= 40 * deltaTime

		foregroundPosition.X -= 80 * deltaTime
		foreground2Position.X -= 80 * deltaTime

		if backgroundPosition.X <= -float32(backgroundTexture.Width)*2 {
			backgroundPosition.X = 0.0
			background2Position.X = float32(backgroundTexture.Width) * 2
		}

		if midgroundPosition.X <= -float32(midgroundTexture.Width)*2 {
			midgroundPosition.X = 0.0
			midground2Position.X = float32(midgroundTexture.Width) * 2
		}

		if foregroundPosition.X <= -float32(foregroundTexture.Width)*2 {
			foregroundPosition.X = 0.0
			foreground2Position.X = float32(foregroundTexture.Width) * 2
		}

		rl.DrawTextureEx(backgroundTexture, backgroundPosition, 0, 2, rl.White)
		rl.DrawTextureEx(backgroundTexture, background2Position, 0, 2, rl.White)

		rl.DrawTextureEx(midgroundTexture, midgroundPosition, 0, 2, rl.White)
		rl.DrawTextureEx(midgroundTexture, midground2Position, 0, 2, rl.White)

		rl.DrawTextureEx(foregroundTexture, foregroundPosition, 0, 2, rl.White)
		rl.DrawTextureEx(foregroundTexture, foreground2Position, 0, 2, rl.White)

		if isOnGround(scarfy, windowHeight) {
			scarfy.velocity = 0
			scarfyIsInAir = false
		} else {
			scarfy.velocity += gravity * deltaTime
			scarfyIsInAir = true
		}

		if rl.IsKeyPressed(rl.KeySpace) && !scarfyIsInAir {
			scarfy.velocity = jumpVelocity
		}

		for i := range enemiesNebula {
			enemiesNebula[i].pos.X += enemiesNebula[i].velocity * deltaTime
		}

		scarfy.pos.Y += scarfy.velocity * deltaTime

		if !scarfyIsInAir {
			scarfy.updateAnimData(deltaTime)
		}

		for i := range enemiesNebula {
			enemy := &enemiesNebula[i]

			nebulaRec := rl.NewRectangle(
				enemy.pos.X+safetyPad,
				enemy.pos.Y+safetyPad,
				enemy.rec.Width-2*safetyPad,
				enemy.rec.Height-2*safetyPad,
			)

			scarfyRec := rl.NewRectangle(
				scarfy.pos.X,
				scarfy.pos.Y,
				scarfy.rec.Width,
				scarfy.rec.Height,
			)

			if rl.CheckCollisionRecs(nebulaRec, scarfyRec) {
				collision = true
			}
		}

		finishLine += enemiesNebula[0].velocity * deltaTime

		if collision {
			rl.DrawText("Game Over!", windowWidth/4, windowHeight/2, 40, rl.White)
		} else {
			for i := range enemiesNebula {
				enemy := &enemiesNebula[i]
				enemy.updateAnimData(deltaTime)

				rl.DrawTextureRec(nebulaTexture, enemy.rec, enemy.pos, rl.White)
			}

			rl.DrawTextureRec(scarfyTexture, scarfy.rec, scarfy.pos, rl.White)
		}

		if !collision && scarfy.pos.X >= finishLine+60 {
			rl.DrawText("You Win!", windowWidth/4, windowHeight/2, 40, rl.White)
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(backgroundTexture)
	rl.UnloadTexture(midgroundTexture)
	rl.UnloadTexture(foregroundTexture)

	rl.UnloadTexture(scarfyTexture)
	rl.UnloadTexture(nebulaTexture)

	rl.CloseWindow()
}
