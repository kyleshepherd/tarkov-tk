module.exports = async (msg) => {
	let helpMsg = '**Thanks for using Tarkov TK** \n \n';
	helpMsg += 'Tarkov TK is very much a work in progress, so if you have any suggestions or issues, please let me know via Twitter https://twitter.com/KyleShepherdDev \n\n';
	helpMsg += 'Also if you enjoy the bot and want to support me, any help would be appreciated! https://www.patreon.com/tarkovtk';
	await msg.channel.send(helpMsg);
};
