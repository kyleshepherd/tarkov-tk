const firebase = require('firebase');

module.exports = async (msg) => {
	firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).get()
		.then((killsQuery) => {
			let kills = [];
			kills = killsQuery.docs.map(kill => ({
				killer: kill.data().killer,
			}));
			let players = [];
			players = kills.reduce((unique, o) => {
				if(!unique.some(obj => obj === o.killer)) {
					unique.push(o.killer);
				}
				return unique;
			},[]);

			let playerKills = [];
			for (const player of players) {
				const playerKillCount = kills.filter(kill => kill.killer === player).length;
				playerKills.push({player: player, killCount: playerKillCount});
			}
			playerKills = playerKills.sort((a,b) => b.killCount - a.killCount);


			let killMsg = '**Most Team Kills\n**';
			for (let i = 0; i < playerKills.length; i++) {
				killMsg += (i + 1) + '. **' + msg.guild.member(playerKills[i].player).nickname + '** - ' + playerKills[i].killCount + ' TKs\n';
			}

			msg.channel.send(killMsg);
		})
		.catch((error) => {
			console.log('Error getting documents: ', error);
		});
};
