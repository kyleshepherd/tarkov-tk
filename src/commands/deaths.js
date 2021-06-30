const firebase = require('firebase');

module.exports = async (msg) => {
	firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).get()
		.then((killsQuery) => {
			let kills = [];
			kills = killsQuery.docs.map(kill => ({
				victim: kill.data().victim,
			}));
			let players = [];
			players = kills.reduce((unique, o) => {
				if(!unique.some(obj => obj === o.victim)) {
					unique.push(o.victim);
				}
				return unique;
			},[]);

			let playerDeaths = [];
			for (const player of players) {
				const playerDeathCount = kills.filter(kill => kill.victim === player).length;
				playerDeaths.push({player: player, killCount: playerDeathCount});
			}
			playerDeaths = playerDeaths.sort((a,b) => b.killCount - a.killCount);


			let killMsg = '**Most Team Deaths\n**';
			for (let i = 0; i < playerDeaths.length; i++) {
				killMsg += (i + 1) + '. **' + msg.guild.member(playerDeaths[i].player).nickname + '** - ' + playerDeaths[i].killCount + ' TDs\n';
			}

			msg.channel.send(killMsg);
		})
		.catch((error) => {
			console.log('Error getting documents: ', error);
		});
};
