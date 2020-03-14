var db = require('../db_helper');

module.exports = async (msg) => {
	var players = await get_players();
	for (let i = 0; i < players.length; i++) {
		let playerKills = await get_player_kills(players[i].player_id);
		players[i].kills = playerKills[0].killCount;
	}
	players = sortByKey(players, 'kills');
	
	let killMsg = '**Most Team Kills \n**';
	for (let i = 0; i < players.length; i++) {
		killMsg += (i + 1) + '. <@' + players[i].player_id + '> - ' + players[i].kills + ' TKs \n';
	}
	await msg.channel.send(killMsg);
};

function get_players()
{
	return new Promise((resolve, reject) => {
		var sql = 'SELECT player_id FROM players;';
		db.query(sql, function (err, result) {
			if (err) reject(err);
			resolve(result);
		});
	});
}

function get_player_kills(player_id) {
	return new Promise((resolve, reject) => {
		var sql = 'SELECT COUNT(*) AS killCount FROM kills WHERE killer = ?;';
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