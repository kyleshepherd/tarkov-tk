var db = require('./db_helper');

module.exports = async (msg) => {
	return new Promise((resolve) => {
		let exists = false;
		var playersResult = check_player_table(msg);
		playersResult.then(function(result) {
			if (result) {
				exists = true;
			} else {
				exists = false;
			}
			resolve(exists);
		});
	});
};

function check_player_table(msg) {
	return new Promise((resolve) => {
		var sql = 'SELECT 1 from players_' + msg.guild.id + ' LIMIT 1;';
		db.query(sql, function (err, result) {
			resolve(result);
		});
	});
}