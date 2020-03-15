var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg, args) => {
	var existing = db_checker(msg);
	existing.then(async function(result) {
		if (result) {
			if (args.length < 2) {
				await msg.channel.send('Make sure you tag 2 users in !tklog \n e.g. `!tklog @Killer @Victim`');
			} else {
				if (msg.mentions.users.size < 2) {
					await msg.channel.send('Make sure you tag 2 users in !tklog \n e.g. `!tklog @Killer @Victim`');
				} else {
					const iterator = msg.mentions.users.values();
		
					const killer = iterator.next().value;
					const victim = iterator.next().value;
					var date = formatDate(new Date());
		
					checkForPlayer(killer, msg);
					checkForPlayer(victim, msg);

					if (args.length > 2) {
						var reasonArr = args.splice(2);
						var reason = reasonArr.join(' ');
					}
		
					var killLog = '';
					if (reason) {
						killLog = 'INSERT INTO kills_' + msg.guild.id + ' (killer, victim, reason, date) VALUES ("' + killer.id + '", "' + victim.id + '", "' + reason + '", "' + date + '");';

					} else {
						killLog = 'INSERT INTO kills_' + msg.guild.id + ' (killer, victim, date) VALUES ("' + killer.id + '", "' + victim.id + '", "' + date + '");';
					}
					db.query(killLog, async function (err) {
						if (err) throw err;
						await msg.channel.send('Kill by ' + killer.username + ' on ' + victim.username + ' logged.');
					});
				}
			}
		} else {
			await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
		}
	});
	
};

function checkForPlayer(player, msg) {
	const playerCheck = 'SELECT * FROM players_' + msg.guild.id + ' WHERE player_id = ' + player.id + ';';
	db.query(playerCheck, function (err, result) {
		if (err) throw err;
		if (result === undefined || result.length == 0) {
			const insertPlayer = 'INSERT INTO players_' + msg.guild.id + ' (player_id, name) VALUES ("' + player.id + '", "' + player.username + '");';
			db.query(insertPlayer);
		} 
	});
}

function formatDate(date) {
	var d = new Date(date),
		month = '' + (d.getMonth() + 1),
		day = '' + d.getDate(),
		year = d.getFullYear();

	if (month.length < 2) 
		month = '0' + month;
	if (day.length < 2) 
		day = '0' + day;

	return [year, month, day].join('-');
}