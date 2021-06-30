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
			const date = formatDate(new Date());

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
				date: date,
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

function formatDate(date) {
	let d = new Date(date),
		month = '' + (d.getMonth() + 1),
		day = '' + d.getDate(),
		year = d.getFullYear();

	if (month.length < 2)
		month = '0' + month;
	if (day.length < 2)
		day = '0' + day;

	return [year, month, day].join('-');
}
