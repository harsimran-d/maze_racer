import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:maze_racer/maze/maze_controller.dart';

class Player {
  Player({required this.id, required this.row, required this.column});
  final String id;
  final int row;
  final int column;
}

class PlayersController extends Notifier<Map<String, Player>> {
  @override
  Map<String, Player> build() {
    return {};
  }

  void move(String move) {
    print("sending move");
    ref.read(mazeProvider.notifier).sendMove(move);
  }

  void add(dynamic rawData) async {
    final id = rawData["playerId"] as String;
    final row = rawData["position"]["row"] as int;
    final column = rawData["position"]["column"] as int;
    state[id] = Player(id: id, row: row, column: column);
    state = {...state};
  }
}

final playersProvider =
    NotifierProvider<PlayersController, Map<String, Player>>(() {
  return PlayersController();
});
