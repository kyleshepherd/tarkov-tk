const Discord = require("discord.js");
const firebase = require("firebase/app");
require("firebase/firestore");

const commandHandler = require("./commands");

require("dotenv").config();

const client = new Discord.Client();
const disbut = require("discord-buttons");
disbut(client);

const firebaseConfig = {
  apiKey: process.env.FIREBASE_API_KEY,
  authDomain: process.env.FIREBASE_AUTH_DOMAIN,
  projectId: process.env.FIREBASE_PROJECT_ID,
  storageBucket: process.env.FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.FIREBASE_APP_ID,
};
// Initialize Firebase
firebase.initializeApp(firebaseConfig);

client.on("message", commandHandler);

client.on("clickButton", async (button) => {
  await firebase
    .firestore()
    .collection("kills")
    .doc(button.id)
    .delete()
    .then(() => {
      button.reply.send("Kill removed");
      button.message.delete();
    })
    .catch((error) => {
      console.error("Error removing kill: ", error);
    });
});

client.login(process.env.BOT_TOKEN);
