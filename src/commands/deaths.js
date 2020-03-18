var db = require('../db_helper');
const db_checker = require('../db_checker');

module.exports = async (msg) => {
	var existing = db_checker(msg);
	existing.then(async function(result) {
		if (result) {
			var players = await get_players(msg);
			for (let i = 0; i < players.length; i++) {
				let playerDeaths = await get_player_deaths(players[i].player_id, msg);
				players[i].deaths = playerDeaths[0].deathCount;
			}
			players = sortByKey(players, 'deaths');
			
			let deathMsg = '**Most Team Deaths \n**';
			for (let i = 0; i < players.length; i++) {
				var playerName = await get_player_name(msg, players[i].player_id);
				deathMsg += (i + 1) + '. **' + playerName + '** - ' + players[i].deaths + ' TDs \n';
			}
			await msg.channel.send(deathMsg);
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

function get_player_deaths(player_id, msg) {
	return new Promise((resolve, reject) => {
		var sql = 'SELECT COUNT(*) AS deathCount FROM kills_' + msg.guild.id + ' WHERE victim = ?;';
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