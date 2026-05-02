# Scope — MVP and beyond

## MVP (target: hackathon demo, 1–3 days of dev)

### In

- Catalog of **15+ curated recipes**, seeded from version-controlled files.
- **Swipe discovery** with three reactions: ❤️ Like, ✕ Dislike, ⭐ Try later.
- **Search by ingredients** (2–5 ingredients, case- and accent-insensitive).
- **Smart filters**:
  - Search across all recipes excluding disliked.
  - Search across liked recipes only.
- **My recipes**: list of `LIKE` and `TRY_LATER` reactions.
- **Cooking School (lite)**: 3–4 static Markdown articles served by API.
- **Auth**: Firebase ID token verification (mobile already integrated).
- **Persistence**: PostgreSQL.

### Out

- ❌ Community photos and voting.
- ❌ User-submitted recipes.
- ❌ Weekly meal plans.
- ❌ Anti-waste suggestions from leftovers.
- ❌ Quizzes / progression in Cooking School.
- ❌ Family/collaborative mode.

### Definition of Done — MVP

- A clean clone runs the full demo with `make up && make migrate && make seed && make run`.
- All endpoints listed in `architecture/api.md` work and are integration-tested.
- 15+ recipes and ≥ 3 lessons are seeded.
- Median API response time on `/api/recipes/discover` < 100 ms locally.
- Swipe → filter by ingredients → My recipes scenario can be demoed end-to-end.

## V1 (post-hackathon)

- Community photos with upvote/downvote, top photo wins.
- User-submitted recipes (with moderation).
- Weekly meal-plan generator from preferences.
- User profile with dietary tags / allergies.

## V2

- Anti-waste suggestions (recipes that maximize use of provided leftovers).
- Cooking School v2: quizzes, challenges, progression tracking.
- Collaborative family mode for shared meal planning.

## Explicit non-goals (forever)

- Ads, sponsored placements, affiliate links.
- Real-time LLM-generated recipes shown to end users.
- Long-form blog content embedded in recipes.
