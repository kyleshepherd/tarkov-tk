# Privacy Policy — Tarkov TK

**Last updated: 27 August 2026**

Tarkov TK ("the bot") is a Discord bot that lets a Discord server log and view
"team kills" for the game Escape From Tarkov. This policy explains what data the
bot stores, why, and how to have it deleted.

The bot is operated by Kyle Shepherd as a personal, non-commercial project.

## What the bot stores

When a member uses the `/tklog` command, the bot stores one record containing:

- The **Discord server (guild) ID** the command was used in
- The **Discord user ID** of the person recorded as the killer
- The **Discord user ID** of the person recorded as the victim
- An optional **free-text reason** typed by the user running the command (max 500 characters)
- The **date and time** the record was created

That is the entirety of what is written to storage.

## What the bot does not store

The bot does **not** store, and does not have access to:

- Message content. The bot does not use the Message Content intent and cannot
  read your messages. It only receives the slash commands explicitly sent to it.
- Usernames, nicknames, display names, avatars, or profile information
- Email addresses, IP addresses, or payment information
- Presence or activity data. The bot does not use the Presence intent.
- Voice data
- Direct messages

## Server member information

To display leaderboards and statistics, the bot needs to turn the stored Discord
user IDs back into readable names. To do this it uses the Server Members
(Guild Members) privileged intent to look up the current members of the server at
the moment a command is run.

This member information is held **in memory only**, for the duration of producing
that single response, and is then discarded. It is never written to the database.
This means names shown by the bot always reflect a member's current nickname or
username, and a member who leaves the server no longer resolves to a name.

## Why this data is stored

The stored records exist solely to provide the bot's features: the team kill log,
the `/tkkills` and `/tkdeaths` leaderboards, and the `/tkstats` statistics and CSV
export. The data is not used for any other purpose.

## Where the data is stored

Records are stored in Google Cloud Firestore, operated by Google. Google acts as a
data processor for this storage. See the
[Google Cloud Privacy Notice](https://cloud.google.com/terms/cloud-privacy-notice)
for details of their practices.

## Sharing

The data is not sold, rented, shared, or transferred to any third party. It is not
used for advertising, profiling, or training machine learning models. The only
third party involved is the storage provider named above.

Data logged in a server is only ever displayed back within that same server.

## Retention and deletion

Records are kept until they are deleted. You can delete data at any time:

- **`/tkremove`** — deletes the most recently logged record in the server.
- **`/tkreset`** — deletes **all** records for the server it is run in. This is
  immediate and permanent.

Removing the bot from a server does not automatically delete that server's records.
If you want data deleted and can no longer run the commands above, contact us using
the details below and we will delete it.

Because the bot stores Discord user IDs, you may request deletion of records
referring to you. Please include your Discord user ID and, where possible, the
server concerned.

## Children

The bot is not directed at children. Discord requires users to be at least 13 years
old, or older where local law requires.

## Changes to this policy

This policy may be updated. Material changes will be reflected by an updated "Last
updated" date at the top of this document, which is published publicly in the
project's repository.

## Contact

Questions or deletion requests:

- Open an issue at https://github.com/kyleshepherd/tarkov-tk/issues
- Or message [@KyleShepherdDev](https://twitter.com/KyleShepherdDev) on Twitter/X
