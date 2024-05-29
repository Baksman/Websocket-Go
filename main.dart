import 'dart:async';
import 'dart:math' as math;
import 'dart:isolate';

// Function to calculate PI using the Monte Carlo method
double calculatePi(int numPoints) {
  var random = math.Random();
  var insideCircle = 0;

  for (var i = 0; i < numPoints; i++) {
    var x = random.nextDouble();
    var y = random.nextDouble();
    if (x * x + y * y <= 1) {
      insideCircle++;
    }
  }

  return (insideCircle / numPoints) * 4;
}

// Function to be executed in a separate isolate
void isolateEntryPoint(List val) {
  ReceivePort receivePort = ReceivePort();

  // receivePort.listen((message) {
  // if (message is int) {
  // Receive the number of points to use for calculation
  // var numPoints = message;
  // Calculate PI
  SendPort sendPort = val[0];
  int number = val[1];
  var pi = calculatePi(number);
  sendPort.send(pi);
  // Send the result back to the main isolate

  // }
  // });

  // Send the receive port to the main isolate
  sendPort.send(receivePort.sendPort);
  Isolate.exit();
}

// Future<double> calculatePiInIsolate(int numPoints) async {
//   ReceivePort receivePort = ReceivePort();
//   // Create a new isolate and send the receive port to it
//   await Isolate.spawn(isolateEntryPoint, receivePort.sendPort);

//   // Receive the send port from the new isolate
//   SendPort sendPort = await receivePort.first;

//   // Send the number of points to the new isolate
//   sendPort.send(numPoints);

//   // Receive the calculated value of PI from the new isolate
//   Completer<double> completer = Completer();
//   receivePort.listen((message) {
//     if (message is double) {
//       completer.complete(message);
//       receivePort.close();
//     }
//   });

//   return completer.future;
// }

void main() async {
  // Number of points to use for calculation
  int numPoints = 10000000;
  // Calculate PI in a separate isolate
  // double pi = await calculatePiInIsolate(numPoints);
  ReceivePort receivePort = ReceivePort();
  Completer<double> completer = Completer();
  await Isolate.spawn(isolateEntryPoint, [receivePort.sendPort, 10000]);

  receivePort.listen((message) {
    print('Calculated PI: $message');
    completer.complete(message);
    receivePort.close();
  });

  // print('Calculated PI: $pi');
}
