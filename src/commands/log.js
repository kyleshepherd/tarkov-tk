const firebase = require('firebase');

module.exports = async (msg, args) => {
	if (args.length < 2) {
		await msg.channel.send('Make sure you tag 2 users in !tklog \n e.g. `!tklog @Killer @Victim`');
	} else {
		if (msg.mentions.users.size < 2) {
			await msg.channel.send('Make sure you tag 2 users in !tklog \n e.g. `!tklog @Killer @Victim`');
		} else {
			const iterator = msg.mentions.users.values();

			const killer = iterator.next().value;
			const victim = iterator.next().value;

			let reason = '';

			if (args.length > 2) {
				const reasonArr = args.splice(2);
				reason = reasonArr.join(' ');
			}

			firebase.firestore().collection('kills').add({
				serverId: msg.guild.id,
				killer: killer.id,
				victim: victim.id,
				reason: reason,
				date: new Date(),
			})
				.then(() => {
					msg.channel.send('Kill by **' + killer.username + '** on **' + victim.username + '** logged.');
				})
				.catch((error) => {
					console.error('Error writing document: ', error);
				});
		}
	}
};
