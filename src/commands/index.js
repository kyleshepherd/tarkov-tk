require('dotenv').config();

var client = require('../db_helper');

module.exports = async (msg) => {
	if (msg.guild.id === process.env.SERVER_ID) {
		const message = msg.content.toLowerCase();
		const args = message.split(' ');
		if (args.length === 0 || args[0].charAt(0) !== '!') return;
		console.log(args);
	}
};