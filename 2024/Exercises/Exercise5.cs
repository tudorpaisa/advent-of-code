using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise5 : IExercise
{
    public int GetDay() => 5;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        (var mapping, var index) = BuildPageOrderMapping(input);

        var updates = BuildUpdates(input, index);

        var goodUpdates = updates.Where(u => UpdateIsOk(mapping, u)).ToArray();
        var midPages = goodUpdates.Select(u => GetMiddlePage(u));

        return new Result(goodUpdates.Select(u => GetMiddlePage(u)).Sum() , false);
    }

    private int GetMiddlePage(List<string> update)
    {
        var count = update.Count();
        var mid = count / 2;

        int.TryParse(update[mid], out var midNumber);
        return midNumber;
    }

    private bool UpdateIsOk(Dictionary<string, List<string>> mapping, List<string> update)
    {
        for (var i = 0; i < update.Count(); i++)
        {
            // NOTE: Maybe check if the mapping contains a?
            var a = update[i];
            for (var j = i + 1; j < update.Count(); j++)
            {
                var b = update[j];

                if (mapping.ContainsKey(b) && mapping[b].Contains(a))
                {
                    return false;
                }
            }
        }

        return true;
    }

    private (Dictionary<string, List<string>>, int) BuildPageOrderMapping(string[] input)
    {
        Dictionary<string, List<string>> mapping = new();

        var idx = 0;
        var line = input[idx];

        while (line != "")
        {
            var split = line.Split("|");
            var collection = mapping.GetValueOrDefault(split.First(), new());
            collection.Add(split.Last());
            mapping[split.First()] = collection;

            idx++;
            line = input[idx];
        }
        idx++;

        return ( mapping, idx );
    }

    private List<List<string>> BuildUpdates(string[] input, int idx)
    {
        List<List<string>> updates = new();
        for (var i  = idx; i < input.Count(); i++)
        {
            updates.Add(input[i].Split(",").ToList());
        }
        return updates;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        (var mapping, var index) = BuildPageOrderMapping(input);

        var updates = BuildUpdates(input, index);

        var badUpdates = updates.Where(u => !UpdateIsOk(mapping, u)).ToArray();

        var sortedUpdates = badUpdates.Select(bd => TopSort(bd, mapping));

        return new Result(sortedUpdates.Select(su => GetMiddlePage(su)).Sum(), false);
    }

    private List<string> TopSort(List<string> a, Dictionary<string, List<string>> mapping)
    {
        Dictionary<string, int> counts = a.ToDictionary(i => i, i => 0);
        foreach (var i in a)
        {
            foreach (var j in mapping.GetValueOrDefault(i, []))
            {
                if (counts.ContainsKey(j)) counts[j]++;
            }
        }

        var searchSpace = counts.Where(kvp => kvp.Value == 0 && a.Contains(kvp.Key)).Select(kvp => kvp.Key).ToList();
        List<string> sorted = [];
        while (searchSpace.Count() != 0)
        {
            var v = searchSpace.First();
            searchSpace.RemoveAt(0);
            sorted.Add(v);

            foreach (var dep in mapping.GetValueOrDefault(v, []))
            {
                if (a.Contains(dep))
                {
                    counts[dep]--;
                    if (counts[dep] == 0) searchSpace.Add(dep);

                }
            }
        }
        return sorted;
    }
}
