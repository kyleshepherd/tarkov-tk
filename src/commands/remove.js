const firebase = require('firebase');

module.exports = async (msg) => {
	firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).orderBy('date', 'desc').limit(1).get()
		.then((killsQuery) => {
			killsQuery.forEach((doc) => {
				firebase.firestore().collection('kills').doc(doc.id).delete()
					.catch((error) => {
						console.error('Error removing kill: ', error);
					});
			});
			msg.channel.send('Last logged kill has been removed');
		})
		.catch((error) => {
			console.log('Error getting documents: ', error);
		});
};
