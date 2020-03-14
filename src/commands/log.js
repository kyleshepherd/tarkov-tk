var db = require('../db_helper');

module.exports = async (msg, args) => {
	db.connect(function (err) {
		if (err) throw err;
		console.log('Connected');
	});
	if (args.length < 2) {
		await msg.channel.send('Make sure you tag 2 users in !log \n e.g. `!log @Killer @Victim`');
	} else {
		if (msg.mentions.users.size < 2) {
			await msg.channel.send('Make sure you tag 2 users in !log \n e.g. `!log @Killer @Victim`');
		} else {
			const iterator = msg.mentions.users.values();

			const killer = iterator.next().value;
			const victim = iterator.next().value;

			checkForPlayer(killer, msg);
			checkForPlayer(victim, msg);

			const killLog = 'INSERT INTO kills (killer, victim) VALUES ("' + killer.id + '", "' + victim.id + '");';

			db.query(killLog, async function (err) {
				if (err) throw err;
				await msg.channel.send('Kill by ' + killer.username + ' on ' + victim.username + ' logged.');
			});
		}
	}
	db.end();
};

function checkForPlayer(player, msg) {
	const playerCheck = 'SELECT * FROM players WHERE player_id = ' + player.id + ';';
	db.query(playerCheck, function (err, result) {
		if (err) throw err;
		if (result === undefined || result.length == 0) {
			const insertPlayer = 'INSERT INTO players (player_id, name) VALUES ("' + player.id + '", "' + player.username + '");';
			db.query(insertPlayer);
		} 
	});
}