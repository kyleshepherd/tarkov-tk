const tklog = require('./log');
const tkkills = require('./kills');
const tkdeaths = require('./deaths');
const tkhelp = require('./help');
const tk = require('./stats');
const tkremove = require('./remove');
const tkinfo = require('./info');
const tkreset = require('./reset');

require('dotenv').config();

const commands = {
	tklog,
	tkkills,
	tkdeaths,
	tkhelp,
	tk,
	tkremove,
	tkinfo,
	tkreset
};

module.exports = async (msg) => {
	const message = msg.content;
	const args = message.split(' ');
	if (args.length === 0 || args[0].charAt(0) !== '!') return;
	args[0].toLowerCase();
	const command = args.shift().substr(1);
	if (Object.keys(commands).includes(command)) {
		commands[command](msg, args);
	}
};
