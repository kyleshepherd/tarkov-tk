var db = require('../db_helper');

module.exports = async (msg) => {
	var players = await get_players();
	for (let i = 0; i < players.length; i++) {
		let playerDeaths = await get_player_deaths(players[i].player_id);
		players[i].deaths = playerDeaths[0].deathCount;
	}
	players = sortByKey(players, 'deaths');
	
	let deathMsg = 'Most Team Deaths \n';
	for (let i = 0; i < players.length; i++) {
		deathMsg += (i + 1) + '. <@' + players[i].player_id + '> - ' + players[i].deaths + ' TDs \n';
	}
	await msg.channel.send(deathMsg);
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

function get_player_deaths(player_id) {
	return new Promise((resolve, reject) => {
		var sql = 'SELECT COUNT(*) AS deathCount FROM kills WHERE victim = ?;';
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