var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg) => {
	var existing = db_checker(msg);
	existing.then(async function(result) {
		if (result) {
			var players = await get_players(msg);
			for (let i = 0; i < players.length; i++) {
				let playerKills = await get_player_kills(players[i].player_id, msg);
				players[i].kills = playerKills[0].killCount;
			}
			players = sortByKey(players, 'kills');
			
			let killMsg = '**Most Team Kills \n**';
			for (let i = 0; i < players.length; i++) {
				killMsg += (i + 1) + '. <@' + players[i].player_id + '> - ' + players[i].kills + ' TKs \n';
			}
			await msg.channel.send(killMsg);
		} else {
			await msg.channel.send('Tarkov TK has not been set up on this server. Run `!start` to do so.');
		}
	});
	
};

function get_players(msg)
{
	return new Promise((resolve, reject) => {
		var sql = 'SELECT player_id FROM players_' + msg.guild.id + ';';
		db.query(sql, function (err, result) {
			if (err) reject(err);
			resolve(result);
		});
	});
}

function get_player_kills(player_id, msg) {
	return new Promise((resolve, reject) => {
		var sql = 'SELECT COUNT(*) AS killCount FROM kills_' + msg.guild.id + ' WHERE killer = ?;';
		db.query(sql, player_id, function (err, result) {
			if (err) reject(err);
			resolve(result);
		});
	});
}

function sortByKey(array, key) {
	return array.sort(function(a, b) {
		var x = a[key]; var y = b[key];
		return ((x > y) ? -1 : ((x < y) ? 1 : 0));
	});
}