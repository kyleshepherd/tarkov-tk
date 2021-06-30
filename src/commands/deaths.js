module.exports = async (msg) => {
	// var existing = db_checker(msg);
	// existing.then(async function(result) {
	// 	if (result) {
	// 		var players = await get_players(msg);
	// 		for (let i = 0; i < players.length; i++) {
	// 			let playerDeaths = await get_player_deaths(players[i].player_id, msg);
	// 			players[i].deaths = playerDeaths[0].deathCount;
	// 		}
	// 		players = sortByKey(players, 'deaths');

	// 		let deathMsg = '**Most Team Deaths \n**';
	// 		for (let i = 0; i < players.length; i++) {
	// 			var playerName = await get_player_name(msg, players[i].player_id);
	// 			deathMsg += (i + 1) + '. **' + playerName + '** - ' + players[i].deaths + ' TDs \n';
	// 		}
	// 		await msg.channel.send(deathMsg);
	// 	} else {
	// 		await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
	// 	}
	// });
};
