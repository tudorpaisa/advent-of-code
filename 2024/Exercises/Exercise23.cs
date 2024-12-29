namespace AdventOfCode.Exercises;

public class Exercise23 : IExercise
{
    public int GetDay() => 23;

    private List<List<string>> Parse(string[] input)
    {
        List<List<string>> res = [];
        foreach (var i in input) res.Add(i.Split("-").ToList());
        return res;
    }

    private Dictionary<string, List<string>> CreateAdjMap(List<List<string>> adjList)
    {
        Dictionary<string, List<string>> map = [];
        foreach( var i in adjList )
        {
            var (a, b) = (i[0], i[1]);
            var aList = map.GetValueOrDefault(a, []);
            aList.Add(b);
            map[a] = aList;

            var bList = map.GetValueOrDefault(b, []);
            bList.Add(a);
            map[b] = bList;
        }
        return map;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var nodePairs = Parse(input);
        var map = CreateAdjMap(nodePairs);
        var nodeList = map.Keys.ToList();

        HashSet<List<string>> groupsOfThree = [];
        for (var i = 0; i < nodeList.Count(); i++)
        {
            for (var j = i + 1; j < nodeList.Count(); j++)
            {
                for (var k = j + 1; k < nodeList.Count(); k++)
                {
                    var a = nodeList[i];
                    var b = nodeList[j];
                    var c = nodeList[k];

                    // Console.WriteLine($"{a} {b} {c}");

                    var aOk = map[a].Contains(b) && map[a].Contains(c);
                    var bOk = map[b].Contains(c) && map[b].Contains(a);
                    var cOk = map[c].Contains(b) && map[c].Contains(a);

                    if (aOk && bOk && cOk)
                    {
                        List<string> gr = [a, b, c];
                        gr.Sort();
                        groupsOfThree.Add(gr);
                    }
                }
            }
        }

        return new(groupsOfThree.Where(g => g.Where(n => n.StartsWith("t")).Count() > 0).Count(), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var nodePairs = Parse(input);
        var map = CreateAdjMap(nodePairs);
        var nodeList = map.Keys.Order().ToList();

        List<string> biggestNetwork = [];
        void Visit(List<string> network)
        {
            if (network.Count() > biggestNetwork.Count())
            {
                biggestNetwork = [..network];
            }

            var tails = network.Count() == 0 ? nodeList : nodeList.SkipWhile(c => c.CompareTo(network.Last()) <= 0).ToList();
            foreach (var i in tails)
            {
                if (network.All(n => map[n].Contains(i)))
                {
                    Visit([..network, i]);
                }
            }
        }
        Visit([]);
        biggestNetwork.Sort();

        return new Result(string.Join(",", biggestNetwork), true);
    }
}
