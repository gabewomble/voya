import { component$ } from "@builder.io/qwik";
import { useTripData } from "./layout";

export const MembersTable = component$(() => {
  const { members } = useTripData().value;
  return (
    <table class="table">
      <caption>Trip Members</caption>
      <thead>
        <tr>
          <th>Name</th>
          <th>Email</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {members.map((member) => (
          <tr key={member.id}>
            <td>{member.name}</td>
            <td>{member.email}</td>
            <td>{member.member_status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
});
