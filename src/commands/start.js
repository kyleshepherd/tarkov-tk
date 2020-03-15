var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg) => {
	var existing = db_checker(msg);
	existing.then(function(result) {
		if (result) {
			msg.channel.send('Tarkov TK has already been setup on this server. Use `!tkhelp` to see how to use Tarkov TK.');
		} else {
			const createKillTable = 'CREATE TABLE kills_' + msg.guild.id + ' (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, killer VARCHAR(255), victim VARCHAR(255), value INT(255), rating INT(1));';
			db.query(createKillTable, function (err) {
				if (err) throw err;
				const createPlayerTable = 'CREATE TABLE players_' + msg.guild.id + ' (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, player_id VARCHAR(255), name VARCHAR(255));';
				db.query(createPlayerTable, async function (err) {
					if (err) throw err;
					await msg.channel.send('Tarkov TK is set up and ready to use! Type `!tkhelp` to see what it can do.');
				});
			});
		}
	});
};