const start = require('./start');

require('dotenv').config();

const commands = {
	start,
};

module.exports = async (msg) => {
	if (msg.guild.id === process.env.SERVER_ID) {
		const message = msg.content.toLowerCase();
		const args = message.split(' ');
		if (args.length === 0 || args[0].charAt(0) !== '!') return;
		const command = args.shift().substr(1);
		if (Object.keys(commands).includes(command)) {
			commands[command](msg, args);
		}
	}
};