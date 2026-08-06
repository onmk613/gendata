package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type Emoji struct {
	Emoji          string   `json:"emoji"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version"`
	IOSVersion     string   `json:"ios_version"`
	SentencesEmoji string   `json:"sentences_emoji"`
}

var Emojis []Emoji

func init() {
	if err := json.Unmarshal([]byte(EmojiJsons), &Emojis); err != nil {
		panic(fmt.Sprintf("failed to load emoji data: %v", err))
	}
	if len(Emojis) == 0 {
		panic("emoji data is empty")
	}
}

var SentencesEmoji = []string{
	// Micro / reactions (emoji varies in position)
	"Nice {emoji}!",
	"Mood {emoji}.",
	"Same {emoji}.",
	"Oof {emoji}.",
	"Yikes {emoji}.",
	"Winning! {emoji}",
	"Noted {emoji}.",
	"Done {emoji}.",
	"Verified {emoji}.",
	"Approved {emoji}.",

	// Short combos (start/middle/end mixed)
	"Good morning {emoji}.",
	"{emoji} Good night.",
	"Let’s go {emoji}!",
	"Nailed it {emoji}{emoji}.",
	"Ship it {emoji}.",
	"One more? {emoji}",
	"Be right back {emoji}.",
	"On my way {emoji}.",
	"Almost there… {emoji}",
	"Thank you! {emoji}",

	// Sequences with light prose
	"{emoji} Big win today.",
	"Coffee first {emoji} {emoji}.",
	"Focus mode {emoji}.",
	"Teamwork! {emoji} {emoji}",
	"Weekend vibes {emoji}.",
	"Deep work now {emoji}.",
	"Daily standup time {emoji}.",
	"Heads down {emoji}.",
	"Let’s eat {emoji}.",
	"Back online {emoji}.",

	// Declarative statements
	"Small steps, steady pace {emoji}.",
	"Clean code, clear minds {emoji}.",
	"Reduced noise, increased signal {emoji}.",
	"Fewer bugs, happier users {emoji}.",
	"Results landed in the dashboard {emoji}.",

	// Imperatives
	"Keep it simple {emoji}.",
	"Please review the PR {emoji}.",
	"Take a quick break {emoji}.",
	"Document the change {emoji}.",
	"Test before you ship {emoji}.",
	"Name things clearly {emoji}.",
	"Measure what matters {emoji}.",
	"Mind the latency budget {emoji}.",
	"Protect the happy path {emoji}.",
	"Celebrate small wins! {emoji}",

	// Questions
	"Ready to roll {emoji}?",
	"Any blockers {emoji}?",
	"Does this scale {emoji}?",
	"Who owns this {emoji}?",
	"What’s next {emoji}?",
	"Can we simplify {emoji}?",
	"Is this necessary {emoji}?",
	"What did we learn {emoji}?",
	"Are we aligned {emoji}?",
	"Time for lunch {emoji}?",

	// Status / standup style
	"Done: tests, docs, and cleanup {emoji}.",
	"Reminder: write the changelog {emoji}.",

	// Food / break
	"Coffee break! {emoji}",
	"Hydrate and stretch {emoji}.",
	"Quick snack, back soon {emoji}.",
	"Lunch run—brb {emoji}.",
	"Treat yourself today {emoji}.",
	"Tea time, then tasks {emoji}.",
	"Dinner after the deploy {emoji}.",
	"Dessert to celebrate! {emoji}",
	"Refuel and refocus {emoji}.",
	"Brunch plans confirmed {emoji}.",

	// Fitness / wellness
	"Breathe in, breathe out {emoji}.",
	"Walk and think {emoji}.",
	"Posture check {emoji}!",
	"Micro break—eyes off screen {emoji}.",
	"Stand, stretch, reset {emoji}.",
	"Quick workout complete {emoji}.",
	"Calm mind, sharp code {emoji}.",
	"Water break—now {emoji}.",
	"Sleep early tonight {emoji}.",
	"Wellness first, always {emoji}.",

	// Celebration / gratitude
	"Thank you, team! {emoji}",
	"Amazing work today {emoji}!",
	"Proud of this release {emoji}.",
	"You crushed it {emoji}!",
	"Confetti for the crew! {emoji}{emoji}",
	"Another milestone reached {emoji}.",
	"Great feedback from users {emoji}.",
	"High-five across time zones {emoji}!",
	"Shipped and shining {emoji}.",
	"Champagne later {emoji}?",

	// Caution / alerts
	"Heads-up: incidents possible {emoji}.",
	"Careful with that change {emoji}.",
	"Review the diff twice {emoji}.",
	"Watch your step here {emoji}.",
	"Rate limit enabled {emoji}.",
	"Safeguards in place {emoji}.",
	"Proceed with caution {emoji}.",
	"Rollback plan ready {emoji}.",
	"Pager is quiet—for now {emoji}.",
	"Triage starts now {emoji}.",

	// Travel / logistics
	"Boarding soon {emoji}.",
	"Wheels up {emoji}.",
	"Landing shortly {emoji}.",
	"Gate changed—move! {emoji}",
	"Delayed, will update {emoji}.",
	"Taxi to the hotel {emoji}.",
	"Home safe {emoji}.",

	// Tiny stories (emoji sprinkled)
	"Started small, grew fast {emoji}.",
	"Found the bug, fixed it {emoji}.",
	"Drew the map, took the path {emoji}.",
	"Built, tested, shipped {emoji}.",
	"Asked, listened, learned {emoji}.",
	"Slowed down to speed up {emoji}.",
	"Cut scope, kept quality {emoji}.",
	"Weathered the storm {emoji}{emoji}.",
	"Saved the day {emoji}.",
	"Onward {emoji}.",

	// Labels / tags with emoji anchors
	"Priority: High {emoji}.",
	"Status: In progress {emoji}.",

	// Mixed fun
	"Plot twist: it worked {emoji}.",
	"Calm is a superpower {emoji}.",
	"Curiosity unlocked {emoji}.",
	"Details matter {emoji}.",
	"Less, but better {emoji}.",
	"Momentum beats motivation {emoji}.",
	"Consistency compounds {emoji}.",
	"Kindness scales {emoji}.",
	"Taste takes time {emoji}.",
	"Simplicity wins {emoji}.",
}

// Data is pull from https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json
//
//go:embed emoji_data.json
var EmojiJsons string
