var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg) => {
	var existing = db_checker(msg);
	existing.then(async function(result) {
		if (result) {
			
		} else {
			await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
		}
	});
};