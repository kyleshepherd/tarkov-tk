var db = require('../db_helper');

module.exports = (msg) => {
	const tableCheck = 'SELECT 1 from kills LIMIT 1;';
	db.query(tableCheck, async function (err, result) {
		if (result) {
			//Table exists
			await msg.channel.send('Tarkov TK has already been setup on this server. Use `!help` to see how to use Tarkov TK.');
		} else {
			//Table doesn't exist
			const createKillTable = 'CREATE TABLE kills (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, killer VARCHAR(255), victim VARCHAR(255), value INT(255), rating INT(1));';
			db.query(createKillTable, function (err) {
				if (err) throw err;
				const createPlayerTable = 'CREATE TABLE players (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, player_id VARCHAR(255), name VARCHAR(255));';
				db.query(createPlayerTable, async function (err) {
					if (err) throw err;
					await msg.channel.send('Tarkov TK is set up and ready to use! Type `!help` to see what it can do.');
				});
			});
		}
	});
};