var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg, args) => {
	var existing = db_checker(msg);
	existing.then(async function(result) {
		if (result) {
			if (args.length < 1) {
				var sql = 'SELECT * FROM kills_' + msg.guild.id + ' ORDER BY id DESC';
				db.query(sql, async function(err, result) {
					if (err) throw err;
					var statMsg = '**Server Kill Stats** \n \n';
					for (let i = 0; i < result.length; i++) {
						var killerName = await get_player_name(msg, result[i].killer);
						var victimName = await get_player_name(msg, result[i].victim);
						var date = new Date(result[i].date);
						date = date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear();
						statMsg += date + ' - Killer: **' + killerName + '** - Victim: **' + victimName + '** ';
						if (result[i].reason != null) {
							statMsg += '- Reason: "' + result[i].reason + '"';
						}
						statMsg += '\n \n';
					}
					await msg.channel.send(statMsg);
				});
			} else {
				if (msg.mentions.users.size < 1) {
					await msg.channel.send('Make sure you tag a user to see their stats e.g. `!log @Player`');
				} else {
					const iterator = msg.mentions.users.values();
					const player = iterator.next().value;

					const playerCheck = 'SELECT * FROM players_' + msg.guild.id + ' WHERE player_id = ' + player.id + ';';
					db.query(playerCheck, async function (err, result) {
						if (err) throw err;
						if (result === undefined || result.length == 0) {
							await msg.channel.send('The tagged player does not exist in the database.');
						} else {
							const killQuery = 'SELECT * FROM kills_' + msg.guild.id + ' WHERE killer = ' + player.id + ' ORDER BY id DESC;';
							db.query(killQuery, async function (err, result) {
								if (err) throw err;
								if (result === undefined || result.length == 0) {
									await msg.channel.send('<@' + player.id + '> hasn\'t team killed anyone...yet');
								} else {
									var playerName = await get_player_name(msg, player.id)
									var statMsg = '**' + playerName + ' Team Kills:** \n \n';
									for (let i = 0; i < result.length; i++) {
										var victimName = await get_player_name(msg, result[i].victim);
										var date = new Date(result[i].date);
										date = date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear();
										statMsg += date + ' - Victim: **' + victimName + '** ';
										if (result[i].reason != null) {
											statMsg += '- Reason: "' + result[i].reason + '"';
										}
										statMsg += '\n \n';
									}
									await msg.channel.send(statMsg);
								}
							});
						}
					});
				}
			}
		} else {
			await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
		}
	});
};

function get_player_name(msg, player_id) {
	return new Promise((resolve) => {
		var playerObj = msg.client.users.fetch(player_id);
		playerObj.then(function (result) {
			resolve(result.username);
		});
	});
}