module.exports = async (msg) => {
	if (msg.guild.id === process.env.SERVER_ID) {
		console.log(msg.content);
	}
}