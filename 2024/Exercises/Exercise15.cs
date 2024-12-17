namespace AdventOfCode.Exercises;

public class Exercise15 : IExercise
{
    public int GetDay() => 15;

    private class Coords
    {
        public int X { get; set; }
        public int Y { get; set; }
    }

    private Dictionary<char, (int Y, int X)> directions = new()
    {
        { '>', ( 0,  1) },
        { '<', ( 0, -1) },
        { '^', (-1,  0) },
        { 'v', ( 1,  0) },
    };

    private bool IsHorizontal(char c) => "<>".Contains(c);
    private bool IsVertical(char c) => "^v".Contains(c);

    private IMapObject? GetAdjacentMovable(List<List<IMapObject?>> map, int x, int y)
    {
        if (map[y][x] != null && map[y][x].Movable) return map[y][x];
        return null;
    }

    private interface IMapObject
    {
        bool Movable { get; }
        Coords Position { get; set; }
        int Width { get; }

        Coords MoveBy(List<List<IMapObject?>> map, int x, int y);
        void Clear(List<List<IMapObject?>> map);
        void Insert(List<List<IMapObject?>> map);
    }

    private class Wall : IMapObject
    {
        public int Width => 1;
        public bool Movable => false;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y) => Position;
        public void Clear(List<List<IMapObject?>> map) {}
        public void Insert(List<List<IMapObject?>> map) {}

    }

    private class Robot : IMapObject
    {
        public int Width => 1;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (map[newY][newX] != null) return Position;

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
        }
    }

    private class Box : IMapObject
    {
        public int Width => 1;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (map[newY][newX] != null) return Position;

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
        }
    }

    private class WideBox : IMapObject
    {
        public int Width => 2;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (x < 0)
            {
                if (map[newY][newX] != null) return Position;
            }
            else
            {
                if (map[newY][newX] != null || map[newY][newX + 1] != null) return Position;
            }

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
            var adjIt = map[Position.Y][Position.X+1];
            if (adjIt.Position.X != Position.X && adjIt.Position.Y != Position.Y)
            {
                throw new Exception("WTF");
            }
            map[Position.Y][Position.X+1] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
            var adjIt = map[Position.Y][Position.X+1];
            if (adjIt != null)
            {
                throw new Exception($"WTF NOT NULL {adjIt.Position.X} {adjIt.Position.Y} {adjIt.GetType().ToString()}");
            }
            map[Position.Y][Position.X+1] = this;
        }
    }

    private (List<List<IMapObject?>>, List<Coords>, Robot, List<IMapObject>) Parse(string[] input)
    {
        List<List<IMapObject?>> parsed = new();
        List<Coords> directions = new();
        List<IMapObject> boxes = new();
        Robot robot = new();

        bool finishedMap = false;

        for (var i = 0; i < input.Count(); i++)
        {
            if (!finishedMap)
            {
                if (input[i] == "")
                {
                    finishedMap = true;
                    continue;
                }

                List<IMapObject?> parsedRow = new();
                for (var j = 0; j < input[0].Count(); j++)
                {
                    var c = input[i][j];

                    switch (c)
                    {
                        case '#':
                            parsedRow.Add(new Wall { Position = new() { X = j, Y = i } });
                            break;
                        case 'O':
                            var box = new Box { Position = new() { X = j, Y = i } };
                            boxes.Add(box);
                            parsedRow.Add(box);
                            break;
                        case '[':
                            var wideBox = new WideBox { Position = new() { X = j, Y = i } };
                            boxes.Add(wideBox);
                            parsedRow.AddRange([wideBox, wideBox]);
                            // parsedRow.Add(wideBox);
                            // parsedRow.Add(wideBox);
                            j++;
                            break;
                        // case ']':
                        //     break;
                        case '@':
                            robot = new Robot { Position = new() { X = j, Y = i } };
                            parsedRow.Add(robot);
                            break;
                        default:
                            parsedRow.Add(null);
                            break;
                    }
                }
                parsed.Add(parsedRow);
            }
            else
            {
                foreach (var c in input[i])
                {
                    switch (c)
                    {
                        case '<':
                            directions.Add(new() { X = -1, Y = 0 });
                            break;
                        case '^':
                            directions.Add(new() { X = 0, Y = -1 });
                            break;
                        case '>':
                            directions.Add(new() { X = 1, Y = 0 });
                            break;
                        case 'v':
                            directions.Add(new() { X = 0, Y = 1 });
                            break;
                        case 'V':
                            directions.Add(new() { X = 0, Y = 1 });
                            break;
                        default:
                            break;
                    }
                }
            }
        }
        return (parsed, directions, robot, boxes);
    }

    private string[] ExpandMap(string[] input)
    {
        List<string> expandedMap = new();
        var finishedMap = false;
        foreach (var row in input)
        {
            if (!finishedMap)
            {
                if (row == "")
                {
                    finishedMap = true;
                    expandedMap.Add("");
                    continue;
                }
                var newRow = "";
                foreach (var c in row)
                {
                    switch (c)
                    {
                        case '#':
                            newRow = newRow + "##";
                            break;
                        case '@':
                            newRow = newRow + "@.";
                            break;
                        case '.':
                            newRow = newRow + "..";
                            break;
                        case 'O':
                            newRow = newRow + "[]";
                            break;
                        default:
                            break;
                    }
                }
                expandedMap.Add(newRow);
            }
            else
            {
                expandedMap.Add(row);
            }
        }
        return expandedMap.ToArray();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        (var map, var directions, var robot, var boxes) = Parse(input);

        foreach (var dir in directions)
        {
            List<IMapObject> itemsToMove = [robot];
            List<IMapObject> moveableChecks = [robot];

            while (moveableChecks.Count() > 0)
            {
                var item = moveableChecks.First();
                moveableChecks.RemoveAt(0);

                var nextX = item.Position.X + dir.X;
                var nextY = item.Position.Y + dir.Y;

                var next = GetAdjacentMovable(map, nextX, nextY);
                if (next != null)
                {
                    moveableChecks.Add(next);
                    itemsToMove = itemsToMove.Prepend(next).ToList();
                }
            }

            foreach (var i in itemsToMove)
            {
                i.Clear(map);
                i.MoveBy(map, dir.X, dir.Y);
                i.Insert(map);
            }
        }

        return new Result(boxes.Select(b => 100 * b.Position.Y + b.Position.X).Sum(), true);
    }

    private (int Y, int X) GetRobotCoords(List<List<string>> map)
    {
        for (var i = 0; i < map.Count(); i++)
        {
            for (var j = 0; j < map[0].Count(); j++)
            {
                if ( map[i][j] == "@" )
                {
                    return (i, j);
                }
            }
        }
        return (-1, -1);
    }

    private (string[] Map, string Movements) SplitMapDirections(string[] input)
    {
        for (var i = 0; i < input.Count(); i++)
        {
            if (input[i] == "")
            {
                return (input[0..i], string.Join("", input[(i+1)..input.Count()]));
            }
        }
        return (new string[] {}, "");
    }

    private (int Y, int X) FindRestOfBox(List<List<string>> map, int Y, int X)
    {
        if (map[Y][X] == "[") return (Y, X + 1);
        return (Y, X - 1);
    }

    private List<(int Y, int X)> FindBoxes(List<List<string>> map)
    {
        List<(int Y, int X)> boxes = new();

        for (var i = 0; i < map.Count(); i++)
        {
            for (var j = 0; j < map[0].Count(); j++)
            {
                if (map[i][j] == "[")
                {
                    boxes.Add((i, j));
                }
            }
        }
        return boxes;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        (var smallMap, var movements) = SplitMapDirections(input);
        var map = ExpandMap(smallMap).Select(i => i.Select(c => c.ToString()).ToList()).ToList();
        (var y, var x) = GetRobotCoords(map);
        map[y][x] = ".";

        foreach (var i in movements)
        {
            var dir = directions[i];
            var newX = x + dir.X;
            var newY = y + dir.Y;

            if ( map[newY][newX] == "." )
            {
                x = newX;
                y = newY;
            }
            else if ( map[newY][newX] == "#" ) continue;
            else
            {
                if (IsHorizontal(i))
                {
                    HashSet<int> marked = new();

                    var searchX = newX;
                    while (true)
                    {
                        searchX += dir.X;
                        if (map[newY][searchX] == "#")
                        {
                            marked = new();
                            break;
                        }
                        else if (map[newY][searchX] == ".")
                        {
                            map[newY][searchX] = i == '<' ? "[" : "]";
                            x += dir.X;
                            y += dir.Y;
                            map[y][x] = ".";
                            break;
                        }
                        else if (map[newY][searchX] == "[" || map[newY][searchX] == "]")
                        {
                            marked.Add(searchX);
                        }
                    }

                    if (marked.Count() == 0) continue;

                    var markedList = marked.ToList();
                    var xVal = markedList.First();
                    map[y][xVal] = i == '>' ? "[" : "]";
                    for (var idx = 1; idx < markedList.Count(); idx++)
                    {
                        var v = markedList[idx];
                        map[y][v] = map[y][xVal] == "]" ? "[" : "]";
                        xVal = v;
                    }
                }
                if (IsVertical(i))
                {
                    List<(int Y, int X)> connected = new();
                    connected.Add((newY, newX));
                    connected.Add(FindRestOfBox(map, newY, newX));

                    var it = 0;
                    while (it < connected.Count())
                    {
                        var connBox = connected[it];
                        // if (map[connBox.Y + dir.Y][connBox.X] == "#")
                        // {
                        //     connected = [];
                        //     break;
                        // }
                        if (map[connBox.Y + dir.Y][connBox.X] == "[" || map[connBox.Y + dir.Y][connBox.X] == "]")
                        {
                            connected.AddRange([(connBox.Y + dir.Y, connBox.X), FindRestOfBox(map, connBox.Y + dir.Y, connBox.X)]);
                        }
                        it++;
                    }
                    if (connected.Where(t => map[t.Y + dir.Y][t.X] == "#").Any())
                    {
                        continue;
                    }
                    connected = connected.Distinct().ToList();
                    connected.Sort();
                    if (dir.Y >= 0) connected.Reverse();
                    foreach (var c in connected)
                    {
                        map[c.Y + dir.Y][c.X] = map[c.Y][c.X];
                        map[c.Y][c.X] = ".";
                    }
                    x = newX;
                    y = newY;
                }
            }
        }

        var allBoxes = FindBoxes(map);

        return new Result(allBoxes.Select(c => c.Y * 100 + c.X).Sum(), false);
    }
}
