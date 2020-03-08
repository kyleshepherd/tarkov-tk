const Discord = require('discord.js');

const commandHandler = require('./commands');

require('dotenv').config();

const client = new Discord.Client();

client.once('ready', async() => {
	console.log('Ready!');
});

client.on('message', commandHandler);

client.login(process.env.BOT_TOKEN);