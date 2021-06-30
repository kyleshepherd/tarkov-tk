const firebase = require('firebase');

module.exports = async (msg) => {
	firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).get()
		.then((killsQuery) => {
			killsQuery.forEach((doc) => {
				firebase.firestore().collection('kills').doc(doc.id).delete()
					.catch((error) => {
						console.error('Error resetting server: ', error);
					});
			});
			msg.channel.send('Your TK server data has been reset');
		})
		.catch((error) => {
			console.log('Error getting documents: ', error);
		});
};
