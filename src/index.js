const Discord = require('discord.js');

const commandHandler = require('./commands');

require('dotenv').config();

const client = new Discord.Client();

client.once('ready', async() => {
	console.log('Ready!');
	// var count = 0;
	// client.guilds.cache.forEach(() => {
	// 	count++;
	// });
	// console.log(count);
});

client.on('message', commandHandler);

client.login(process.env.BOT_TOKEN);