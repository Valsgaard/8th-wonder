# API improvements, suggestions and musings

## User creation

Okay this is a though one, because from the documentation I'm reading the intent is to have a "new user" on every session, but at the same time you want to have a "friend list" feature.
If you want users to be persistent you need to add in user account, and not just rely on a username alone.

So user creation is now a unique username, and a password. (Or OAuth2 in any shape, like Facebook, Google, etc)

## Game state

Preface: Access pattern doesn't dictate data model, but using a pattern different from the model is just unnecessary confusion.

First, currently the game state is represented as a part of the user: `/user/<userid>/state`

This puts the game state as a child to a user, locking out the possibility of multiple users belonging to the same "game", in one shape or another.
I'd put the game state as an entity of its own, to remove this restriction from it. Opening up the possibilities of multiplayer options, even if they aren't actually within the scope of the game. Things change, let's make sure the model is open for iteration.

Secondly, the game state is being updated completely by the client. Instead of adding in the score and gamesPlayed count, just add a score for each game and let the backend check if that score is an improvement, and increment the game play count.
(This doesn't change the fact that the user can upload a random value, but that won't change unless you run the game on the server)

## Friends

Sending the entire friend list each time isn't feasible, send a "addFriend" request with a single userID at a time, and have a "removeFriend" endpoint as well to remove them.

Now "Friends" in and of itself right now, doesn't serve any real purpose other than having cherry picked high score board. Which is kinda an odd thing to have, since there's no global high score board in the first place, so why a friend specific one? I assume that there's more to it.

[Back to README](README.MD)
