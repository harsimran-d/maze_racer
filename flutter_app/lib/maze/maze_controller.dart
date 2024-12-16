import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:maze_racer/maze/maze.dart';
import 'package:maze_racer/maze/player_position_provider.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class MazeController extends AutoDisposeNotifier<Maze?> {
  MazeController();
  late WebSocketChannel _socket;

  @override
  Maze? build() {
    ref.onDispose(() {
      print("disposing ");
      _socket.sink.close();
    });
    getMazeFromOnline();
    return null;
  }

  void sendMove(String move) {
    _socket.sink.add(move);
  }

  void getMazeFromOnline() async {
    _socket = WebSocketChannel.connect(Uri.parse("ws://localhost:3000/game"));
    await _socket.ready;
    _socket.stream.listen((data) {
      final rawData = jsonDecode(data);
      print(rawData);
      if (rawData["MazeMap"] != null) {
        List<List<int>> finalList = [];
        final list = rawData["MazeMap"] as List;

        for (final child in list) {
          List<int> innerList = [];
          for (final val in child) {
            innerList.add(val as int);
          }
          finalList.add(innerList);
        }
        state = Maze(mazeMap: finalList, size: rawData["Size"] as int);
      } else if (rawData["playerId"] != null) {
        ref.read(playersProvider.notifier).add(rawData);
      }
    });
  }
}

final mazeProvider = AutoDisposeNotifierProvider<MazeController, Maze?>(
  () {
    return MazeController();
  },
);
