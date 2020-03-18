const tkstart = require('./start');
const tklog = require('./log');
const tkkills = require('./kills');
const tkdeaths = require('./deaths');
const tkhelp = require('./help');
const tk = require('./stats');
const tkremove = require('./remove');

require('dotenv').config();

const commands = {
	tkstart,
	tklog,
	tkkills,
	tkdeaths,
	tkhelp,
	tk,
	tkremove
};

module.exports = async (msg) => {
	const message = msg.content.toLowerCase();
	const args = message.split(' ');
	if (args.length === 0 || args[0].charAt(0) !== '!') return;
	const command = args.shift().substr(1);
	if (Object.keys(commands).includes(command)) {
		commands[command](msg, args);
	}
};