package controllers

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

// 获取短信内容的URL
type GetMessageContentURLController struct {
	BaseController
}

func (c *GetMessageContentURLController) Get() {
	serviceType := c.Ctx.Input.Param(":type")
	serviceIndexNum := c.Ctx.Input.Param(":index")
	if serviceType == "game" {
		if serviceIndexNum == "" {
			result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(GameList))))
			messageContentURL := "http://game.eldogame.com/content/game/" + result.String()
			c.StringResult(messageContentURL)
		} else {
			indexInt, _ := strconv.Atoi(serviceIndexNum)
			c.RedirectURL(GameList[indexInt])
		}
	}
	c.StringResult("404")
}

var GameList = []string{
	"https://d3huw0u63gtszr.cloudfront.net/en/wasteland-warriors/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/hamster-roll/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/school-bus-pickup/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/confident-driver/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/millionaire-quiz/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/grenade-toss/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/soccer-pro/index.html",
	"http://d3huw0u63gtszr.cloudfront.net/en/defend-the-beach/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/giant-hamster-run/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/bob-and-chainsaw/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/puzzle-ball/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/rainbow-stacker/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/christmas-tree-fun/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/true-love-calculator/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/picture-slider-animals/index.html",
	"http://d3huw0u63gtszr.cloudfront.net/en/chef-slash/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/baby-cow-launcher/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/world-of-words/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/pet-hop/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/super-boxing/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/javelin-olympics/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/jetpack-blast/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/candy-jam/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/street-fight/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/guess-the-soccer-star/index.html",
	"http://s3-ap-southeast-1.amazonaws.com/marketjs-shoal1/en/slam-dunk-forever/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/stack-the-burger/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/nerd-quiz/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/mahjong-pyramids/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/tictactoe-with-friends/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/monster-mahjong/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/word-hunter/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/dear-grim-reaper/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/zero-collsion/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/robots-vs-aliens/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/pizza-cafe/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/crossthebridge/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/3-card-tarot-reading/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/spider-solitaire/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/office-fight/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/mountain-hop/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/santa-city-run/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/spot-the-difference/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/connect-the-dots/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/car-factory/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/candy-timberman/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/clash-of-vikings/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/slice-of-zen/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/social-blackjack/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/math-genius/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/bubble-pop-adventures/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/alien-claw-crane/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/boat-dash/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/taxi-pickup/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/space-chasers/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/tank-battle/index.html",
	"http://d3huw0u63gtszr.cloudfront.net/en/donut-slam-dunk/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/hostage-rescue/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/montezuma-gems/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/basketball-legend/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/santa-delivery/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/casual-checkers/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/red-outpost/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/nine-pool-game/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/count-faster/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/desert-rally/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/bingo-world/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/unblock-it/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/lost-in-time/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/cubeform/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/monsters-and-cake/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/trump-quiz-game/index.html",
	"http://d3huw0u63gtszr.cloudfront.net/en/4-in-a-row/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/rabbit-zombie-defense/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/car-park-puzzle/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/adventure-craft/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/sudoku-village/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/wildwestshootout/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/touchdown-pro/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/new-team-kaboom/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/candy-slide/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/vampires-and-garlic/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/blackjack-vegas-21/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/templecrossing/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/penalty-kick-game/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/zombie-survival/index.html",
	"https://d3huw0u63gtszr.cloudfront.net/en/sheep-jump/index.html",
}
